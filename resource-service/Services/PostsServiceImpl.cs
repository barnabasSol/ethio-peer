using Grpc.Core;
using ResourceService.Grpc.ProtoBuf;
using ResourceService.Repositories;

namespace ResourceService.Services;

public class PostsServiceImpl : PostsService.PostsServiceBase
{
    private readonly PostRepo _postsRepository;

    public PostsServiceImpl(PostRepo postsRepository)
    {
        _postsRepository = postsRepository;
    }
    public override async Task<PostsReply> GetPosts(PostsRequest request,ServerCallContext context)
    {
        var posts = await _postsRepository.GetPostsByRoomId(Guid.Parse(request.RoomId));
        var reply = new PostsReply();
        reply.Posts.AddRange(posts.Select(p => new Post
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