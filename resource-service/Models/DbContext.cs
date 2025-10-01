using Microsoft.EntityFrameworkCore;
namespace ResourceService.Models;


public class Context : DbContext
{
    public Context() { }
    public Context(DbContextOptions<Context> options)
            : base(options) { }
    protected override void OnModelCreating(ModelBuilder modelBuilder)
    {
        modelBuilder.Entity<Topic>()
                       .HasKey(b => b.TopicId);
        modelBuilder.Entity<Course>()
       .HasKey(b => b.CourseCode);
        modelBuilder.Entity<Document>()
       .HasKey(b => b.Id);
    }
    public DbSet<Course> Courses { get; set; }
    public DbSet<Topic> Topics { get; set; }
    public DbSet<Document> Documents { get; set; }
    public DbSet<Room> Rooms { get; set; }
    public DbSet<Post> Posts { get; set; }
    public DbSet<RoomMember> RoomMembers { get; set; }
}

