namespace ResourceService.Models.Dtos
{
    public class RoomDTO
    {
        public string SessionId { get; set; }= string.Empty;
        public Guid TopicId { get; set; }

        public string Name { get; set; } = string.Empty;
    }
}