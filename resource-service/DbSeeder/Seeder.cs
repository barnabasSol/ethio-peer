using System.Text.Json;
using ResourceService.Models;

public static class DataSeeder
{
    public static async Task SeedCoursesAsync(Context context)
    {

        var jsonData = await File.ReadAllTextAsync("./DbSeeder/courses.json");
        var courses = JsonSerializer.Deserialize<List<Course>>(jsonData, new JsonSerializerOptions
        {
            PropertyNameCaseInsensitive = true
        });

        if (courses != null)
        {
            context.Courses.AddRange(courses);
            await context.SaveChangesAsync();
        }
    }
    public static async Task SeedTopicsAsync(Context context)
    {
        
        var jsonData = await File.ReadAllTextAsync("./DbSeeder/topics.json");
        var topics = JsonSerializer.Deserialize<List<Topic>>(jsonData, new JsonSerializerOptions
        {
            PropertyNameCaseInsensitive = true
        });

        if (topics != null)
        {
            context.Topics.AddRange(topics);
            await context.SaveChangesAsync();
        }
    }
}
