using Grpc.Core;
using ResourceService.Grpc.ProtoBuf;
using ResourceService.Models;
using ResourceService.Repositories;

namespace ResourceService.Services;

public class PostsServiceImpl : PostsService.PostsServiceBase
{
    private readonly PostRepo _postsRepository;

    public PostsServiceImpl(PostRepo postsRepository)
    {
        _postsRepository = postsRepository;
    }
    public override async Task<PostsReply> GetPosts(PostsRequest request, ServerCallContext context)
    {
        var pagedQuery = new PagedQuery
    {
        PageSize = request.PageSize,
        PageNumber = request.PageNumber
    };
        var posts = await _postsRepository.GetPostsByRoomId(Guid.Parse(request.RoomId),pagedQuery); 
        var reply=new PostsReply();
        reply.Posts.AddRange(posts.Items.Select(p => new Grpc.ProtoBuf.Post
        {
            Id = p.Id.ToString(),
            SenderId = p.SenderId.ToString(),
            Content = p.Content,
            RoomId = p.RoomId.ToString(),
            CreatedAt = p.CreatedAt.ToString("o") // ISO 8601 format
        }));
        return reply; 
    }
}