using Consul;
using gateway;
using gateway.Service;
using Microsoft.AspNetCore.Authentication.JwtBearer;
using Microsoft.IdentityModel.Tokens;

var builder = WebApplication.CreateBuilder(args);

builder.Services.AddSingleton<IJwksProvider>(
    new JwksProvider("http://localhost:2000/static/.well-known/jwks.json")
);

builder.Services.AddHostedService<JwksRefreshService>();

builder
    .Services.AddAuthentication(options =>
    {
        options.DefaultAuthenticateScheme = JwtBearerDefaults.AuthenticationScheme;
        options.DefaultChallengeScheme = JwtBearerDefaults.AuthenticationScheme;
    })
    .AddJwtBearer();

builder
    .Services.AddOptions<JwtBearerOptions>(JwtBearerDefaults.AuthenticationScheme)
    .Configure<IJwksProvider>(
        (options, jwksProvider) =>
        {
            options.RequireHttpsMetadata = false;

            options.TokenValidationParameters = new TokenValidationParameters
            {
                ValidateIssuer = true,
                ValidIssuer = "ep.barney-host.site",

                ValidateAudience = true,
                ValidAudiences = ["ep-web.barney-host.site"],

                ValidateLifetime = true,
                RequireExpirationTime = true,

                ValidateIssuerSigningKey = true,

                IssuerSigningKeyResolver = (token, securityToken, kid, validationParameters) =>
                {
                    return jwksProvider.GetSigningKeys().Where(k => k.KeyId == kid) ?? [];
                },
            };
        }
    );

builder
    .Services.AddAuthorizationBuilder()
    .AddPolicy(
        "authenticated",
        policy =>
        {
            policy.RequireAuthenticatedUser();
        }
    );

builder.Services.AddCors(options =>
{
    options.AddPolicy(
        "WebOriginCorsPolicy",
        policy =>
        {
            policy
                .WithOrigins("http://localhost:5173")
                .AllowAnyHeader()
                .AllowAnyMethod()
                .AllowCredentials();
        }
    );
});

builder.Services.AddEndpointsApiExplorer();
builder.Services.AddHealthChecks();
builder.Services.AddSwaggerGen();

builder.Services.AddControllers();

builder
    .Services.AddReverseProxy()
    .LoadFromMemory(gateway.Routes.GetRoutes(), Clusters.GetClusters());

var app = builder.Build();

if (app.Environment.IsDevelopment())
{
    app.UseSwagger();
    app.UseSwaggerUI();
}

app.MapControllers();
app.MapHealthChecks("/health");

app.UseCors();

app.UseAuthorization();

app.MapReverseProxy();

app.Run();
