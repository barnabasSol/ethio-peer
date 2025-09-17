using ResourceService.Repositories;
using Microsoft.EntityFrameworkCore;
using Minio;
using ResourceService.Models;
using ResourceService.Services;

var builder = WebApplication.CreateBuilder(args);

// Add services to the container.
builder.Services.AddDbContext<Context>(options =>
    options.UseNpgsql(builder.Configuration.GetConnectionString("DefaultConnection")!));
builder.Services.AddSingleton<Rabbit>();
builder.Services.AddScoped<CourseRepo, CourseRepo>();
builder.Services.AddScoped<TopicRepo, TopicRepo>();
builder.Services.AddScoped<DocRepo, DocRepo>();
builder.Services.AddScoped<RoomRepository, RoomRepository>();
builder.Services.AddScoped<PostRepo, PostRepo>();
builder.Services.AddControllers();
builder.Services.AddOpenApi();
builder.Services.AddEndpointsApiExplorer();
builder.Services.AddSwaggerGen();
builder.Services.AddSignalR();
Console.WriteLine("Minio endpoint: " + builder.Configuration["Minio:Endpoint"]);
// Add CORS
builder.Services.AddCors(options =>
{
    options.AddPolicy("AllowVueClient",
        policy =>
        {
            policy.WithOrigins("http://localhost:5173", "http://localhost:5174")
                  .AllowAnyHeader()
                  .AllowAnyMethod()
                  .AllowCredentials(); // needed for SignalR
        });
});
builder.Services.AddMinio(configureClient => configureClient
            .WithEndpoint(builder.Configuration["Minio:Endpoint"])
            .WithCredentials(builder.Configuration["Minio:AccessKey"], builder.Configuration["Minio:SecretKey"])
            .WithSSL(false)
        .Build());
var app = builder.Build();
// using (var scope = app.Services.CreateScope())
// {
//     var db = scope.ServiceProvider.GetRequiredService<Context>();
//     db.Database.Migrate();
// }

// Configure the HTTP request pipeline.
var rabbit = app.Services.GetRequiredService<Rabbit>();
await rabbit.Subscribe();
app.UseSwagger();
app.UseSwaggerUI(c =>
{
    c.SwaggerEndpoint("/swagger/v1/swagger.json", "Course Topic Service API V1");
    c.RoutePrefix = string.Empty; // Set Swagger UI at the app's root
});
app.MapOpenApi();
app.MapHub<RoomHub>("/roomHub");
app.UseCors("AllowVueClient");

// app.UseHttpsRedirection();

app.MapControllers();
app.Lifetime.ApplicationStopping.Register(() =>
{
    var rabbit = app.Services.GetRequiredService<Rabbit>();
    rabbit.Close().Wait();
});
app.Run();


