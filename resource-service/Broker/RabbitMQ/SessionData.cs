namespace ResourceService.Broker.RabbitMQ;

public class SessionData
{
    public Guid SessionId { get; set; }
    public string UserName { get; set; } = string.Empty;
    public Guid TopicId { get; set; }
    
}