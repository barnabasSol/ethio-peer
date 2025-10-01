using System.Text.Json.Serialization;

namespace ResourceService.Models;

public record DocDTO
{
    // public Guid Id { get; set; }
    [JsonIgnore]
    public IFormFile? DocFile { get; set; }
    public Guid UploaderId { get; set; }
    // public string Title { get; set; } = string.Empty;  d
    public Guid RoomId { get; set; }
}