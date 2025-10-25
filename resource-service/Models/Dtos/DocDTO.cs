using System.Text.Json.Serialization;

namespace ResourceService.Models.Dtos;

public record DocDTO
{
    public string Title { get; set; } = string.Empty; 
    // public string FileName { get; set; } = "";
    public string UploaderId { get; set; } 
    public Guid RoomId { get; set; }
}