using Microsoft.EntityFrameworkCore;

namespace ResourceService.Models;

[PrimaryKey("CourseCode")]
public class Course
{
    public required string CourseCode { get; set; }
    public required string Name { get; set; }
    public required string Description { get; set; }
    public int CreditHour { get; set; }
    public CourseCategory Category { get; set; }
}
public enum CourseCategory
{
Programming_and_SoftwareDev,
Systems_and_Infrastructure,
Databases_and_DataMgmt,
Web_and_Mobile_Tech,
Specialized_and_EmergingAreas,
IT_Management_and_Research
}