namespace ResourceService.Models.Dtos;

public class RoomResp
{
    public Guid RoomId { get; set; }
    public string RoomName { get; set; }
    public int MemberCount { get; set; }
    public string TopicName { get; set; }
    public string CourseCode { get; set; }

}