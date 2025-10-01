using System.ComponentModel.DataAnnotations;
using System.ComponentModel.DataAnnotations.Schema;

namespace ResourceService.Models;

public class Topic
{
    [Key]
    public Guid TopicId { get; set; }= Guid.NewGuid();  
    public required string Name { get; set; }
    public required string Description { get; set; }    
    public string CourseCode { get; set; } = String.Empty;
    [ForeignKey("CourseCode")]
    public Course? Course { get; set; }
    public List<Room> Rooms { get; set; } = [];
}