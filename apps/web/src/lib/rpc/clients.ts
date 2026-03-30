import { createClient } from "@connectrpc/connect"
import { createConnectTransport } from "@connectrpc/connect-web"

import { UserService } from "@/lib/rpc/gen/user/user_pb"
import { AuthService } from "@/lib/rpc/gen/auth/auth_pb"

const transport = createConnectTransport({
    baseUrl: "http://127.0.0.1:8080",
})
export const userClient = createClient(UserService, transport)
export const authClient = createClient(AuthService, transport)
