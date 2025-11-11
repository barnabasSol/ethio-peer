using System.Text.Json.Serialization;

namespace ResourceService.Models
{
    public class Post
    {
        public Guid Id { get; set; }
        public string SenderId { get; set; } = "";
        public string SenderName { get; set; } = "";
        public string Content { get; set; } = "";
        public bool IsDoc { get; set; } = false;
        public string DocKey { get; set; } = "";
        public string DocTitle { get; set; } = "";
        public Guid RoomId { get; set; }
        [JsonIgnore]
        public Room? Room { get; set; }
        public DateTime CreatedAt { get; set; }

    }
}
