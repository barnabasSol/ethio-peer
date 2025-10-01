using System.Text.Json.Serialization;
using Microsoft.EntityFrameworkCore;

namespace ResourceService.Models;

[PrimaryKey("CourseCode")]
public class Course
{
    public required string CourseCode { get; set; }
    public required string Name { get; set; }
    public required string Description { get; set; }
    public int CreditHour { get; set; }
    [JsonIgnore]
    public List<Topic> Topics { get; set; } = [];
    public CourseCategory Category { get; set; }
}
