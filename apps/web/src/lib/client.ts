import { createClient } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";
import { AuthService } from "$lib/gen/auth/auth_pb";
import { UserService } from "$lib/gen/user/user_pb";
import { BlobService } from "$lib/gen/blob/blob_pb";

const transport = createConnectTransport({
    baseUrl: "http://localhost:8080",
});

export const authClient = createClient(AuthService, transport);
export const userClient = createClient(UserService, transport);
export const blobClient = createClient(BlobService, transport);
