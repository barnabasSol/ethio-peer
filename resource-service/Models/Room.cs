using System.Text.Json.Serialization;

namespace ResourceService.Models
{
    public class Room
    {
        public Guid Id { get; set; }
        public string SessionId { get; set; } = string.Empty;
        public Guid TopicId { get; set; }
        [JsonIgnore]
        public Topic? Topic { get; set; }

        public string Name { get; set; } = string.Empty; 
        [JsonIgnore]
        public List<Post> Posts { get; set; } = [];
        [JsonIgnore]
        public List<Document> Documents { get; set; } = [];
        // public List<RoomMember> Members { get; set; } = [];


    
    }
}