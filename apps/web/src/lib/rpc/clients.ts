import { createClient } from "@connectrpc/connect"
import { createConnectTransport } from "@connectrpc/connect-web"

import { UserService } from "@/lib/rpc/gen/user/user_pb"
import { AuthService } from "@/lib/rpc/gen/auth/auth_pb"

const transport = createConnectTransport({
    baseUrl: process.env.NEXT_PUBLIC_API_URL!,
})

export const userClient = createClient(UserService, transport)
export const authClient = createClient(AuthService, transport)
