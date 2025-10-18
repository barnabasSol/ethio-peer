using Microsoft.EntityFrameworkCore;
using Minio;
using ResourceService.Models;
using ResourceService.Models.Dtos;

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
        public async Task<PagedList<Post>> GetPostsByRoomId(Guid roomId, PagedQuery pq)
        {   var room = await _context.Rooms.Where(r => r.Id == roomId).FirstOrDefaultAsync() ?? throw new ArgumentNullException(nameof(roomId), "Room not found");
            var postsQuery = _context.Posts.Where(p => p.RoomId == roomId).OrderByDescending(p => p.CreatedAt).AsQueryable();
        //     var pagedPosts = await PagedList<Post>.CreateAsync(postsQuery, pq.PageSize, pq.PageNumber);
        //     var projectedPosts = pagedPosts.Items.Select(async p =>
        //    {
        //        var resp = new PostResp
        //        {
        //            Id = p.Id,
        //            SenderId = p.SenderId,
        //            Content = p.Content,
        //            IsDoc = p.IsDoc,
        //            DocTitle = p.DocTitle,
        //            RoomId = p.RoomId,
        //            CreatedAt = p.CreatedAt
        //        };
        //        if (p.IsDoc && p.DocKey != "")
        //        {
        //            resp.DocUrl = await _minio.GenerateDownloadLink(p.DocKey);
        //        }
        //        return resp;
        //    });
        //     var postResps = await Task.WhenAll(projectedPosts);

            return await PagedList<Post>.CreateAsync(postsQuery,pq.PageSize,pq.PageNumber);
        }

        public async Task<string> GetWeeklyPosts()
        {
            var oneWeekAgo = DateTime.UtcNow.AddDays(-7);
            var roomPosts = await _context.Posts.Where(p => p.CreatedAt >= oneWeekAgo && p.Room != null).Select(
                p => new
                {
                    RoomName = p.Room!.Name,
                    PostedAt = p.CreatedAt,
                    Content = p.Content
                }
            ).ToListAsync();
            var groupedPosts = roomPosts.GroupBy(r => r.RoomName).Select(
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
            return formattedQuery;

        }
    }
}
