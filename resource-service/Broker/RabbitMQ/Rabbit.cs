using Microsoft.EntityFrameworkCore.Metadata;
using RabbitMQ.Client;
using RabbitMQ.Client.Events;
using ResourceService.Broker.RabbitMQ;
using ResourceService.Broker.RabbitMQ.Dtos;
using ResourceService.Models;
using ResourceService.Models.Dtos;
using ResourceService.Repositories;
using System.Text;
using System.Text.Json;
public class Rabbit
{
    private IChannel? _channel;
    private IConnection? _connection;
    private readonly IConfiguration _configuration;
    private readonly IServiceProvider _serviceProvider;

    public Rabbit(IServiceProvider serviceProvider, IConfiguration config)
    {
        _serviceProvider = serviceProvider;
        _configuration = config;
    }

    public async Task InitiateConsuming()
    {
        string exchangeName = "Session_Exg";
        await ConnectRabbit();
        if (_channel == null)
        {   //what if throwing exception mm
            Console.WriteLine("Channel not initialized");
            return;
        }

        await _channel.ExchangeDeclareAsync(exchange: exchangeName, type: ExchangeType.Topic, durable: true);

        // we create queue with a generated name which we subscribe for listening to session creation:
        QueueDeclareOk queueDeclareResult = await _channel.QueueDeclareAsync();
        string queueName = queueDeclareResult.QueueName;
        await _channel.QueueBindAsync(queue: queueName, exchange: exchangeName, routingKey: "session.#");


        var consumer = new AsyncEventingBasicConsumer(_channel);
        consumer.ReceivedAsync += async (model, ea) =>
        {
            var routingKey = ea.RoutingKey; 
            var byteArray = ea.Body.ToArray();
            var message = Encoding.UTF8.GetString(byteArray);
            try
            {
                switch (routingKey)
                {
                    case "session.created":
                        await ProcessSessionCreation(ea, message);
                        break;

                    case "session.member.joined":
                        await ProcessMemberJoined(ea, message);
                        break;

                    default:
                        Console.WriteLine($"Unhandled routing key: {routingKey}");
                        break;
                }
            }
            catch (Exception ex)
            {
                Console.WriteLine($"Error processing message: {ex.Message}");
                await _channel.BasicNackAsync(ea.DeliveryTag, multiple: false, requeue: false);
            }

        };

        await _channel.BasicConsumeAsync(queueName, autoAck: false, consumer: consumer);
    }

    private async Task ProcessMemberJoined(BasicDeliverEventArgs ea, string message)
    {
        var scope = _serviceProvider.CreateScope();
        var roomRepo = scope.ServiceProvider.GetRequiredService<RoomRepository>();
        var opts = new JsonSerializerOptions { PropertyNameCaseInsensitive = true };
        var member = JsonSerializer.Deserialize<MemberData>(message, opts); 
        if (member != null)
        {
            Console.WriteLine($"Received SessionId={member.SessionId}, MemberId={member.MemberId}");
            var roomId = await roomRepo.GetRoomIdBySessionId(member.SessionId);
            RoomMemberDTO dto = new RoomMemberDTO
            {
                RoomId = roomId,
                MemberId = member.MemberId
            };
            roomRepo.AddParticipant(dto).Wait();
            await _channel!.BasicAckAsync(ea.DeliveryTag, multiple: false);
        }
        else
        {
            Console.WriteLine("Received message could not be deserialized to MemberData.");
            return;
        }



    }

    private async Task ProcessSessionCreation(BasicDeliverEventArgs ea, string message)
    {
        var scope = _serviceProvider.CreateScope();
        var roomRepo = scope.ServiceProvider.GetRequiredService<RoomRepository>();
        var topicRepo = scope.ServiceProvider.GetRequiredService<TopicRepo>();
        var opts = new JsonSerializerOptions { PropertyNameCaseInsensitive = true };
        var session = JsonSerializer.Deserialize<SessionData>(message, opts);

        if (session != null)
        {
            Console.WriteLine($"Received SessionId={session.SessionId}, UserName={session.UserName}, TopicId={session.TopicId}");
            string topicName = await topicRepo.GetTopicNameById(Guid.Parse(session.TopicId));
            RoomDTO roomDto = new RoomDTO
            {
                SessionId =  session.SessionId,
                Name = session.UserName.ToUpper() + "'s " + topicName,
                TopicId = Guid.Parse(session.TopicId)

            };
            var room = await roomRepo.AddRoom(roomDto);
            await roomRepo.AddParticipant(new RoomMemberDTO { RoomId = room.Id, MemberId = session.OwnerId });
            await _channel!.BasicAckAsync(ea.DeliveryTag, multiple: false);

        }
        else
        {
            Console.WriteLine("Received message could not be deserialized to SessionData.");
        }
    }

    private async Task ConnectRabbit()
    {
        if (_connection != null && _channel != null && _channel.IsOpen) return;

        var factory = new ConnectionFactory
        {
            HostName = _configuration["RabbitMQ:Host"]!,
            UserName = _configuration["RabbitMQ:Username"]!,
            Password = _configuration["RabbitMQ:Password"]!
        };
        _connection = await factory.CreateConnectionAsync();
        _channel = await _connection.CreateChannelAsync();
        return;
    }
    public async Task Close()
    {
        try { await _channel?.CloseAsync()!; }
        catch
        {
            Console.WriteLine("Channel could not be closed");
        }
        try { await _connection?.CloseAsync()!; }
        catch
        {
            Console.WriteLine("Connection could not be closed");

        }
    }
}