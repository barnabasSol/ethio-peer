using Microsoft.EntityFrameworkCore;
using Minio;
using Minio.DataModel;
using Minio.DataModel.Args;
using ResourceService.Models;
using ResourceService.Models.Dtos;

namespace ResourceService.Repositories;

public class DocRepo(IMinioClient minioClient, Context context)
{
    private readonly IMinioClient _minioClient = minioClient;
    private readonly Context _context = context;
    private readonly string bucketName = "docs";

    public async Task<DocUploadResp> AddDoc(DocDTO dto)
    {
        try
        {
            var topic = _context.Rooms.Find(dto.RoomId) ?? throw new ArgumentNullException("dto.RoomId", "There is no room with the given id");
            bool found = await _minioClient.BucketExistsAsync(new BucketExistsArgs().WithBucket(bucketName));
            if (!found)
            {
                await _minioClient.MakeBucketAsync(new MakeBucketArgs().WithBucket(bucketName));
            }
            string objectKey = Guid.NewGuid().ToString();
            var url = await _minioClient.PresignedPutObjectAsync(
                            new PresignedPutObjectArgs()
                                .WithBucket(bucketName)
                                .WithObject(objectKey)
                                .WithExpiry(3600));

            Document doc = new()
            {
                Id = new Guid(),
                ObjectKey = objectKey,
                UploaderId = dto.UploaderId,
                Title = dto.Title,
                DateUploaded = DateTime.UtcNow,
                RoomId = dto.RoomId
            };
            await _context.Documents.AddAsync(doc);
            Post post = new()
            {
                Id = Guid.NewGuid(),
                RoomId = dto.RoomId,
                SenderId = dto.UploaderId,
                SenderName = dto.SenderName,
                IsDoc = true,
                DocKey = objectKey,
                DocTitle = dto.Title,
                CreatedAt = DateTime.UtcNow
            };
            await _context.Posts.AddAsync(post);
            await _context.SaveChangesAsync();
            return new DocUploadResp
            {
                UploadUrl = url,
                DocKey = objectKey,
                DocTitle = dto.Title
            };
        }
        catch (FileNotFoundException)
        {
            throw;
        }
        catch (Exception e)
        {
            Console.WriteLine($"The exception source is {e.Source} and the message \n {e.Message}");
            throw;
        }
    }
    public async Task<string> GenerateDownloadLink(string key)
    {
        //check if the file exists
        try
        {
            var bucketExists = await _minioClient.BucketExistsAsync(new BucketExistsArgs().WithBucket(bucketName));
            if (!bucketExists)
            {
                return "";
            }
            string presignedUrl = await _minioClient.PresignedGetObjectAsync(
    new PresignedGetObjectArgs()
        .WithBucket(bucketName)
        .WithObject(key)
        .WithExpiry(3600)
            );
            return presignedUrl;


        }
        catch (Exception e)
        {
            Console.WriteLine("Could not generate a link", e.Message);
            return "";
        }
    }
    public async Task<Document> ModifyDocMetadata(Guid id, DocDTO dto)
    {
        Document? doc = _context.Documents.Find(id) ?? throw new FileNotFoundException();
        ArgumentNullException.ThrowIfNull(dto);
        doc.RoomId = dto.RoomId;
        doc.Title = dto.Title;
        _context.Documents.Update(doc);
        await _context.SaveChangesAsync();
        return doc;
    }
    public async Task<bool> DeleteDoc(Guid id)
    {
        var doc = await _context.Documents.FindAsync(id);
        if (doc == null)
        {
            return false;
        }
        try
        {
            //check if the file exists
            StatObjectArgs statObjectArgs = new StatObjectArgs().WithBucket(bucketName).WithObject(doc.ObjectKey);
            var stat = await _minioClient.StatObjectAsync(statObjectArgs);
        }
        catch (Minio.Exceptions.ObjectNotFoundException)
        {
            throw new FileNotFoundException();
        }
        //if it does, delete the file from minio
        RemoveObjectArgs removeObjectArgs = new RemoveObjectArgs().WithBucket(bucketName).WithObject(doc.ObjectKey);
        await _minioClient.RemoveObjectAsync(removeObjectArgs);
        //delete the file metadata from the database
        _context.Documents.Remove(doc);
        await _context.SaveChangesAsync();
        return true;
    }

    public async Task<List<DocResp>> GetDocsAsync(CourseCategory? category = null, int count = 0)
    {
        List<DocResp> docSugges = [];

        //base case
        var documentsQuery = _context.Documents.Where(d => d.Room != null && d.Room.Topic != null).Include(d => d.Room).ThenInclude(r => r!.Topic)
        .AsNoTracking().AsQueryable();

        if (category != null)
        {
            documentsQuery = documentsQuery.Where(d => d.Room!.Topic!.Course != null && d.Room.Topic.Course.Category == category);
        }
        if (count > 0)
        {
            documentsQuery = documentsQuery.Take(count);

        }
        var docList = await documentsQuery.OrderByDescending(x => x.DateUploaded).ToListAsync();

        foreach (var doc in docList)
        {
            // var link = await GenerateDownloadLink(doc);
            docSugges.Add(new DocResp
            {
                DocId = doc.Id,
                DocTitle = doc.Title,
                TopicName = doc.Room!.Topic!.Name,
                UploadDate = doc.DateUploaded,
                DocKey = doc.ObjectKey

            });
        }

        return [.. docSugges];
    }


    public async Task<List<DocResp>> GetDocumentSuggestionsAsync(List<string> courses)
    {
        if (courses == null || courses.Count == 0)
        {
            return await GetDocsAsync(count: 3);

        }
        var docs = await _context
                    .Documents
                    .AsNoTracking()
                    .Where(d =>
                     d.Room != null
                     && d.Room.Topic != null
                     && d.Room.Topic.Course != null
                     && courses.Contains(d.Room.Topic.Course.Name))
                    .OrderByDescending(d => d.DateUploaded)
                    .Take(3)
                    .Select(d => new DocResp
                    {
                        DocId = d.Id,
                        DocTitle = d.Title,
                        TopicName = d.Room!.Topic!.Name,
                        UploadDate = d.DateUploaded,
                        DocKey = d.ObjectKey
                    })
                    .ToListAsync();
        if (docs == null || docs.Count == 0)
            return await GetDocsAsync(count: 3);
        return docs;
    }
    private string RemoveExtension(string fileName)
    {
        int lastDot = fileName.LastIndexOf('.');
        if (lastDot == -1)
        {
            return fileName;
        }
        return fileName[..fileName.LastIndexOf('.')];
    }

    // public async Task<List<Document>> GetDocsByRoomId(Guid roomId)
    // {

    //     var docs = from doc in _context.Documents
    //                join room in _context.Rooms on doc.RoomId equals room.Id
    //                where doc.RoomId == roomId
    //                select doc;
    //     return await docs.ToListAsync();
    // }
}