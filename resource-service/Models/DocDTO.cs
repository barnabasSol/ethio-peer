using System.Text.Json.Serialization;

namespace ResourceService.Models;

public record DocDTO
{
    public string FileName { get; set; } = "";
    public Guid UploaderId { get; set; } 
    public Guid RoomId { get; set; }
}