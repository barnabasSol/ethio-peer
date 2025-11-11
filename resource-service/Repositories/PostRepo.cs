using Microsoft.EntityFrameworkCore;
using Minio;
using ResourceService.Models;
using ResourceService.Models.Dtos;

namespace ResourceService.Repositories
{
    public class PostRepo(Context cxt, GeminiCaller geminiCaller)
    {
        private readonly Context _context = cxt;
        private readonly GeminiCaller _geminiCaller = geminiCaller;
        //Add Post
        public async Task AddPostedMessage(Guid roomId, string userName, string senderId, string message)
        {
            //check if the room exists
            Room room = await _context
                                .Rooms
                                .FindAsync(roomId) ?? throw new ArgumentNullException("Room not found");
            Post newPost = new()
            {
                Id = Guid.NewGuid(),
                RoomId = roomId,
                SenderId = senderId,
                SenderName = userName,
                Content = message,
                CreatedAt = DateTime.UtcNow
            };
            await _context.Posts.AddAsync(newPost);
            await _context.SaveChangesAsync();

        }

        //Get Posts by RoomId
        public async Task<PagedList<Post>> GetPostsByRoomId(Guid roomId, PagedQuery pq)
        {
            var roomExist = await _context.Rooms
                            .AnyAsync(r => r.Id == roomId);
            if (!roomExist) throw new ArgumentNullException(nameof(roomId), "Room not found");
            var postsQuery = _context.Posts
                                    .AsNoTracking()
                                    .Where(p => p.RoomId == roomId)
                                    .OrderByDescending(p => p.CreatedAt)
                                    .AsQueryable();

            return await PagedList<Post>.CreateAsync(postsQuery, pq.PageSize, pq.PageNumber);
        }

        public async Task<string> GetWeeklyPosts()
        {
            var oneWeekAgo = DateTime.UtcNow.AddDays(-7);
            var roomPosts = await _context
                .Posts
                .Where(p => p.CreatedAt >= oneWeekAgo && p.Room != null)
                .Select(
                    p => new
                    {
                        RoomName = p.Room!.Name,
                        PostedAt = p.CreatedAt,
                        Content = p.Content
                    }
                ).ToListAsync();
            var groupedPosts = roomPosts
            .GroupBy(r => r.RoomName)
            .Select(
                g => new RoomPosts
                {
                    RoomName = g.Key,
                    Posts = g.Select(p => new RoomPosts.PostMinimal
                    {
                        PostedAt = p.PostedAt,
                        Content = p.Content
                    }).ToList()
                }).ToList();
            string formattedQuery = "";
            foreach (var entry in groupedPosts)
            {
                formattedQuery += $"Room: [{entry.RoomName}] -";
                foreach (var post in entry.Posts)
                {
                    formattedQuery += $"User: '{post.Content}' -";
                }
            }
            return await _geminiCaller.CallGeminiApiAsync(formattedQuery);

        }
    }
}
