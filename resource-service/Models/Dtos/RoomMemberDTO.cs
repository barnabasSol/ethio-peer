namespace ResourceService.Models.Dtos;
public class RoomMemberDTO
{
    public Guid RoomId { get; set; }
    public string MemberId { get; set; } = "";
}