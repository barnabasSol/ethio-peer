    using System.Text.Json.Serialization;

namespace ResourceService.Models;

public class CourseDTO
{
    public required string CourseCode { get; set; }
    public required string Name { get; set; }
    public required string Description { get; set; }
    public int CreditHour { get; set; }
    [field: JsonConverter(typeof(JsonStringEnumConverter))] public CourseCategory Category { get; set; } 

}