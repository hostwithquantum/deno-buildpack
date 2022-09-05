import { Server } from 'https://deno.land/std@0.119.0/http/server.ts';

let port = 8080;
if (Deno.env.get('PORT') !== undefined) {
	port = Number(Deno.env.get('PORT'));
}

const handler = (request: Request) => {
	const body = `Hello from Deno!\n\nYour user-agent is:\n\n${
		request.headers.get(
			'user-agent',
		) ?? 'Unknown'
	}`;

	return new Response(body, { status: 200 });
};

const server = new Server({ port, handler });

console.log(`Listening on http://0.0.0.0:${port}`);

await server.listenAndServe();
