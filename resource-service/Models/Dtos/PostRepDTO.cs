using System.Text.Json.Serialization;

namespace ResourceService.Models.Dtos;

public class PostResp
{
        public Guid Id { get; set; }
        public string SenderId { get; set; } = "";
        public string Content { get; set; } = "";
        public bool IsDoc { get; set; } = false;
        public string DocUrl { get; set; } = "";
        public string DocTitle { get; set; } = "";
        public Guid RoomId { get; set; }
        public DateTime CreatedAt { get; set; }
}