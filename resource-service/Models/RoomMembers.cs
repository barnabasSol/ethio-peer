namespace ResourceService.Models;
public class RoomMember
{
    public Guid Id { get; set; } = Guid.NewGuid();
    public Guid RoomId { get; set; }
    public Room? Room { get; set; }
    public string UserId { get; set; } = "";
} 