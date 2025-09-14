using System.Text.Json.Serialization;

namespace ResourceService.Models
{
    public class Room
    {
        public Guid Id { get; set; }
        public Guid SessionId { get; set; }
        public string Name { get; set; } = string.Empty; 
        [JsonIgnore]
        public List<Post> Posts { get; set; } = [];
        [JsonIgnore]
        public List<Document> Documents { get; set; } = [];

    
    }
}