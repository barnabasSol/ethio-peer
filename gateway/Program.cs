var builder = WebApplication.CreateBuilder(args);

builder.Services.AddEndpointsApiExplorer();
builder.Services.AddHealthChecks();
builder.Services.AddSwaggerGen();

builder.Services.AddControllers();

builder.Services.AddReverseProxy().LoadFromMemory([], []);

var app = builder.Build();

if (app.Environment.IsDevelopment())
{
    app.UseSwagger();
    app.UseSwaggerUI();
}

app.MapControllers();
app.MapHealthChecks("/health");

app.UseAuthentication();
app.UseAuthorization();

// app.MapReverseProxy();

app.Run();
