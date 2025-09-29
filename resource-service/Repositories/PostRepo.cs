using Microsoft.EntityFrameworkCore;
using ResourceService.Models;

namespace ResourceService.Repositories
{
    public class PostRepo(Context cxt)
    {
        private readonly Context _context = cxt;
        //Add Post
        public async Task AddPostedMessage(Guid roomId, Guid senderId, string message)
        {
            //check if the room exists
            Room room = await _context.Rooms.FindAsync(roomId) ?? throw new ArgumentNullException("Room not found");
            Post newPost = new()
            {
                Id = Guid.NewGuid(),
                RoomId = roomId,
                SenderId = senderId,
                Content = message,
                CreatedAt = DateTime.UtcNow
            };
            await _context.Posts.AddAsync(newPost);
            await _context.SaveChangesAsync();

        }

        //Get Posts by RoomId
        public async Task<PagedList<Post>> GetPostsByRoomId(Guid roomId,PagedQuery pq)
        {
            var room = await _context.Rooms.Where(r => r.Id == roomId).FirstOrDefaultAsync() ?? throw new ArgumentNullException(nameof(roomId),"Room not found");
            var postsQuery = _context.Posts.Where(p => p.RoomId == roomId).OrderByDescending(p => p.CreatedAt).AsQueryable();
            return await PagedList<Post>.CreateAsync(postsQuery,pq.PageSize,pq.PageNumber);
        }
        
    }
}
