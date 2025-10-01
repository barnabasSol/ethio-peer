using Microsoft.AspNetCore.SignalR;
using ResourceService.Repositories;

namespace ResourceService.Services
{
    public class RoomHub(PostRepo postRepo):Hub
    {
        private readonly PostRepo _postRepo = postRepo;
        public async Task BroadcastToGroup(string roomId,string senderId,string message)
        {
            try
            {
                //persist to db
                _postRepo.AddPostedMessage(Guid.Parse(roomId),Guid.Parse(senderId), message).Wait();
                // Broadcast message to all clients in the specified group (room)
                await Clients.Group(roomId).SendAsync("ReceivedRoomMsg", senderId, message);

            }
            catch (Exception e)
            {
                Console.WriteLine("Error broadcasting message to group"+ e.Message);
            }
        }
        public async Task JoinRoom(Guid roomId)
        {
            //sample Guid roomId = Guid.Parse("d290f1ee-6c54-4b01-90e6-d701748f0851");  
            await Groups.AddToGroupAsync(Context.ConnectionId, roomId.ToString());
        }
    }
}
//# UOW Pattern
//# translation of user id to its name