namespace ResourceService.Models.Dtos;

public class DocResp
{
    public Guid DocId { get; set; }
    public string DocTitle { get; set; } = "";
    public string TopicName { get; set; } = "";
    public DateTime UploadDate { get; set; }
    public string DocKey { get; set; } = "";
    
}