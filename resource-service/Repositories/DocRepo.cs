using Microsoft.EntityFrameworkCore;
using Minio;
using Minio.DataModel;
using Minio.DataModel.Args;
using ResourceService.Models;

namespace ResourceService.Repositories;

public class DocRepo(IMinioClient minioClient, Context context)
{
    private readonly IMinioClient _minioClient = minioClient;
    private readonly Context _context = context; 
    public async Task<DocDTO> AddDoc(DocDTO dto)
    {
        try
        {
            var topic = _context.Rooms.Find(dto.RoomId) ?? throw new ArgumentNullException("dto.RoomId", "There is no room with the given id");
            bool found = await _minioClient.BucketExistsAsync(new BucketExistsArgs().WithBucket(dto.UploaderId.ToString()));
            if (!found)
            {
                await _minioClient.MakeBucketAsync(new MakeBucketArgs().WithBucket(dto.UploaderId.ToString()));
            }
            string fileName = Guid.NewGuid().ToString() + "_" + dto.DocFile!.FileName;
            var putObjectArgs = new PutObjectArgs()
                            .WithBucket(dto.UploaderId.ToString())
                            .WithObject(fileName)
                            .WithStreamData(dto.DocFile.OpenReadStream())
                            .WithObjectSize(dto.DocFile.Length)
                            .WithContentType(dto.DocFile.ContentType)
                            ;
            await _minioClient.PutObjectAsync(putObjectArgs);

            Document doc = new()
            {
                Id = fileName,
                FileName = dto.DocFile.FileName,
                UploaderId = dto.UploaderId,
                // Title = dto.Title, 
                DateUploaded = DateTime.UtcNow,
                RoomId = dto.RoomId
            };
            await _context.Documents.AddAsync(doc);
            await _context.SaveChangesAsync();
            return dto;
        }
        catch (FileNotFoundException)
        {
            throw;
        }
        catch (Exception)
        {
            throw;
        }
    }
    public async Task<Stream> DownloadDoc(string id)
    {

        ObjectStat stat;
        //check if the file exists
        try
        {
            Document? doc = _context.Documents.Find(id) ?? throw new FileNotFoundException();
            StatObjectArgs statObjectArgs = new StatObjectArgs().WithBucket(doc.UploaderId.ToString()).WithObject(id);
            stat = await minioClient.StatObjectAsync(statObjectArgs);
            var memoryStream = new MemoryStream();
            var getObjectArgs = new GetObjectArgs()
                .WithBucket(doc.UploaderId.ToString())
                .WithObject(id)
                .WithCallbackStream(async (stream, _) =>
                {
                    await stream.CopyToAsync(memoryStream);
                });

            await _minioClient.GetObjectAsync(getObjectArgs);
            memoryStream.Position = 0; 
            return memoryStream;
        }
        catch (Minio.Exceptions.ObjectNotFoundException)
        {
            throw new FileNotFoundException();
        }
    }

    public async Task<Document> ModifyDocMetadata(string id, DocDTO dto)
    {
        Document? doc = _context.Documents.Find(id) ?? throw new FileNotFoundException();
        ArgumentNullException.ThrowIfNull(dto); 
        doc.RoomId = dto.RoomId;
        _context.Documents.Update(doc);
        await _context.SaveChangesAsync();
        return doc;
    }
    public async Task<bool> DeleteDoc(string id)
    {
        var doc = await _context.Documents.FindAsync(id);
        if (doc == null)
        {
            return false;
        }
        try
        {
            //check if the file exists
            StatObjectArgs statObjectArgs = new StatObjectArgs().WithBucket(doc.UploaderId.ToString()).WithObject(id);
            var stat = await _minioClient.StatObjectAsync(statObjectArgs);
        }
        catch (Minio.Exceptions.ObjectNotFoundException)
        {
            throw new FileNotFoundException();
        }
        //if it does, delete the file from minio
        RemoveObjectArgs removeObjectArgs = new RemoveObjectArgs().WithBucket(doc.UploaderId.ToString()).WithObject(id);
        await _minioClient.RemoveObjectAsync(removeObjectArgs);
        //delete the file metadata from the database
        _context.Documents.Remove(doc);
        await _context.SaveChangesAsync();
        return true;
    }

    public async Task<List<Document>> GetDocsAsync()
    {
       return await _context.Documents.ToListAsync();
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