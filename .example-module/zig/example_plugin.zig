const std = @import("std");

pub fn range(len: usize) []const u0 {
    return @as([*]u0, undefined)[0..len];
}

export fn lpm_entrypoint(db_path: [*:0]u8, argc: usize, argv: [*][*:0]u8) void {
    std.debug.print("db_path: {s}\n", .{db_path});

    for (range(argc)) |_, i| {
        std.debug.print("db_path: {s}\n", .{argv[i]});
    }
}
