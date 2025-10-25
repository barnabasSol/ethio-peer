using Microsoft.AspNetCore.SignalR;
using ResourceService.Repositories;

namespace ResourceService.Services
{
    public class RoomHub(PostRepo postRepo) : Hub
    {
        private readonly PostRepo _postRepo = postRepo;
        public async Task BroadcastToGroup(string roomId, string senderId,string userName, string? message, string? docKey, string? docTitle)
        {
            try
            {
                if (string.IsNullOrEmpty(docTitle))
                {
                    _postRepo.AddPostedMessage(Guid.Parse(roomId), senderId, message!).Wait();
                    // Broadcast message to all clients in the specified group (room)
                    await Clients.Group(roomId).SendAsync("ReceivedRoomPost",senderId,userName, message);
                    return;
                }
                else
                    await Clients.Group(roomId).SendAsync("ReceivedRoomDoc",senderId, userName, docKey, docTitle);
                return;

            }

            catch (Exception e)
            {
                Console.WriteLine("Error broadcasting message to group" + e.Message);
            }
        }
        public async Task JoinRoom(Guid roomId)
        {  
            await Groups.AddToGroupAsync(Context.ConnectionId, roomId.ToString());
        }
    }
} 