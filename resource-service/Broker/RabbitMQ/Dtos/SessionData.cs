namespace ResourceService.Broker.RabbitMQ.Dtos;

public class SessionData
{
    public Guid SessionId { get; set; }
    public string UserName { get; set; } = string.Empty;
    public Guid OwnerId { get; set; }

    public Guid TopicId { get; set; }

}
