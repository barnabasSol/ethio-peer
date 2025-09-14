using System.Text.Json.Serialization;

namespace ResourceService.Models
{
    public class Post
    {
        public Guid Id { get; set; }
        public Guid SenderId { get; set; }
        public string Content { get; set; } ="";
        public Guid RoomId { get; set; }
        [JsonIgnore]
        public Room? Room { get; set; }
        public DateTime CreatedAt { get; set; }

    }
}
