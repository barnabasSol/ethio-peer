namespace ResourceService.Broker.RabbitMQ.Dtos;

public class SessionData
{
    public string SessionId { get; set; } = string.Empty;
    public string UserName { get; set; } = string.Empty;
    public string OwnerId { get; set; } = string.Empty;
    public string TopicId { get; set; } = string.Empty;


}