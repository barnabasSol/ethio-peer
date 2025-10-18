namespace ResourceService.Models.Dtos;

public class RoomPosts
{
    public string RoomName { get; set; } = string.Empty;
    public List<PostMinimal> Posts { get; set; } = [];
    public class PostMinimal
    {
        public DateTime PostedAt { get; set; }
        public string Content { get; set; } = string.Empty;
}
}