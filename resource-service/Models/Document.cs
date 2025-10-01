using System.Text.Json.Serialization;

namespace ResourceService.Models;

public class Document
{
    public string Id { get; set; } = string.Empty;
    public string FileName { get; set; } = string.Empty;
    public Guid UploaderId { get; set; }
    // public string Title { get; set; } = string.Empty; 
    public DateTime DateUploaded { get; set; }  
    public Guid RoomId { get; set; } // Foreign key to Room

    public Room ?Room { get; set; }  // Navigation property to Room
    //May be add date modified if needed
}