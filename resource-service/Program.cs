using Microsoft.EntityFrameworkCore;
using Minio;
using ResourceService.Models;
using ResourceService.Repositories;
using ResourceService.Services;

var builder = WebApplication.CreateBuilder(args);

builder.Services.AddHealthChecks();

// Add services to the container.
builder.Services.AddDbContext<Context>(options =>
    options.UseNpgsql(builder.Configuration.GetConnectionString("DefaultConnection")!)
);
builder.Services.AddSingleton<Rabbit>();
builder.Services.AddScoped<CourseRepo, CourseRepo>();
builder.Services.AddScoped<TopicRepo, TopicRepo>();
builder.Services.AddScoped<DocRepo, DocRepo>();
builder.Services.AddScoped<RoomRepository, RoomRepository>();
builder.Services.AddScoped<PostRepo, PostRepo>();
builder.Services.AddControllers();
builder.Services.AddEndpointsApiExplorer();
builder.Services.AddSwaggerGen();
builder.Services.AddSignalR();
builder.Services.AddGrpc();
Console.WriteLine("Minio endpoint: " + builder.Configuration["Minio:Endpoint"]);

// Add CORS
builder.Services.AddCors(options =>
{
    options.AddPolicy(
        "AllowVueClient",
        policy =>
        {
            policy
                .WithOrigins(
                    "http://localhost:5173",
                    "http://localhost:5174",
                    "http://127.0.0.1:5500"
                )
                .AllowAnyHeader()
                .AllowAnyMethod()
                .AllowCredentials(); // needed for SignalR
        }
    );
});

builder.Services.AddMinio(configureClient =>
    configureClient
        .WithEndpoint(builder.Configuration["Minio:Endpoint"])
        .WithCredentials(
            builder.Configuration["Minio:AccessKey"],
            builder.Configuration["Minio:SecretKey"]
        )
        .WithSSL(true)
        .Build()
);
var app = builder.Build();

app.MapHealthChecks("/health");

// using (var scope = app.Services.CreateScope())
// {
//     var db = scope.ServiceProvider.GetRequiredService<Context>();
//     db.Database.Migrate();
// }

// Configure the HTTP request pipeline.
var rabbit = app.Services.GetRequiredService<Rabbit>();
await rabbit.InitiateConsuming();
app.UseSwagger();
app.UseSwaggerUI(c =>
{
    c.SwaggerEndpoint("/swagger/v1/swagger.json", "Course Topic Service API V1");
    c.RoutePrefix = string.Empty; // Set Swagger UI at the app's root
});
app.MapHub<RoomHub>("/roomHub");
app.UseCors("AllowVueClient");

// app.UseHttpsRedirection();
app.MapGrpcService<PostsServiceImpl>();
app.MapControllers();

// using (var scope = app.Services.CreateScope())
// {
//     var dbContext = scope.ServiceProvider.GetRequiredService<Context>();
//     await DataSeeder.SeedCoursesAsync(dbContext);
//     await DataSeeder.SeedTopicsAsync(dbContext);
// }
app.Lifetime.ApplicationStopping.Register(() =>
{
    var rabbit = app.Services.GetRequiredService<Rabbit>();
    rabbit.Close().Wait();
});
app.Run();
