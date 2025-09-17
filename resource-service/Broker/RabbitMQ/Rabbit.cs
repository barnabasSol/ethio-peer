using Microsoft.EntityFrameworkCore.Metadata;
using RabbitMQ.Client;
using RabbitMQ.Client.Events;
using ResourceService.Broker.RabbitMQ;
using ResourceService.Models;
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

    public async Task Subscribe()
    {
        await ConnectRabbit();
        if (_channel == null)
        {   //what if throwing exception mm
            Console.WriteLine("Channel not initialized");
            return;
        }

        await _channel.ExchangeDeclareAsync(exchange: "Session_Exg", type: ExchangeType.Topic);

        // we create a non-durable, exclusive, autodelete queue with a generated name:
        QueueDeclareOk queueDeclareResult = await _channel.QueueDeclareAsync();
        string queueName = queueDeclareResult.QueueName;
        await _channel.QueueBindAsync(queue: queueName, exchange: "Session_Exg", routingKey: "session.created");


        var consumer = new AsyncEventingBasicConsumer(_channel);
        consumer.ReceivedAsync += async (model, ea) =>
        {
            //your code to handle the message goes here 
            var byteArray = ea.Body.ToArray();
            var message = Encoding.UTF8.GetString(byteArray);

            try
            {
                using var scope = _serviceProvider.CreateScope();
                var roomRepo = scope.ServiceProvider.GetRequiredService<RoomRepository>();
                var topicRepo = scope.ServiceProvider.GetRequiredService<TopicRepo>();
                var opts = new JsonSerializerOptions { PropertyNameCaseInsensitive = true };
                var session = JsonSerializer.Deserialize<SessionData>(message, opts);

                if (session != null)
                {
                    Console.WriteLine($"Received SessionId={session.SessionId}, UserName={session.UserName}, TopicId={session.TopicId}");
                    string topicName = await topicRepo.GetTopicNameById(session.TopicId);
                    RoomDTO roomDto = new RoomDTO
                    {
                        SessionId = session.SessionId,
                        Name = session.UserName.ToUpper() + "'s " + topicName

                    };
                    await roomRepo.AddRoom(roomDto);
                    await _channel.BasicAckAsync(ea.DeliveryTag, multiple: false);

                }
                else
                {
                    Console.WriteLine("Received message could not be deserialized to SessionData.");
                }
            }
            catch (Exception ex)
            {
                Console.WriteLine($"Error processing message: {ex.Message}");
                await _channel.BasicNackAsync(ea.DeliveryTag, multiple: false, requeue: false);
                // consider negative ack or requeue if using manual acks
            }

        };

        await _channel.BasicConsumeAsync(queueName, autoAck: false, consumer: consumer);
    }

    private async Task ConnectRabbit()
    {
        if (_connection != null && _channel != null && _channel.IsOpen) return;

        var factory = new ConnectionFactory
        {
            HostName = _configuration["RabbitMQ:Host"]!,
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