import { FormEvent, useState } from "react";
import { API_URL } from "./App";
import { useNavigate } from "react-router-dom";

export const SignInPage = () => {
	const [emailOrUsername, setEmailOrUsername] = useState<string>("");
	const [password, setPassword] = useState<string>("");
	const redirect = useNavigate();
	
	const handlerSignIn = async (e: FormEvent) => {
		e.preventDefault();

		const response = await fetch(`${API_URL}/users/signin`, {
			method: "POST",
			body: JSON.stringify({ emailOrUsername, password }),
		});
		if (response.ok) {
			redirect("/");
		} else {
			alert("faild to sign in");
		}
	};
	return (
		<div className="flex bg-background h-screen items-center justify-center p-5">
			<div className="bg-secondary h-[400px] w-[320px] md:h-[550px] md:w-[500px] flex flex-col items-center justify-center rounded-3xl gap-8 text-primary shadow-lg border border-border">
				<h1 className="text-4xl md:text-5xl text-accent font-bold mb-4">
					Sign In
				</h1>
				<form
					action="GET"
					className="flex flex-col w-[300px] md:w-[350px] gap-1 font-semibold"
				>
					<label htmlFor="email-or-username">Email or Username:</label>
					<input
						onChange={(e) => setEmailOrUsername(e.target.value)}
						type="text"
						className="border-1 border-accent px-2 py-2.5 rounded focus:outline-none focus:ring-2 focus:ring-accent mb-6"
						placeholder="email@example.com"
						value={emailOrUsername}
						id="email-or-username"
					/>
					<label htmlFor="password">Password:</label>
					<input
						onChange={(e) => setPassword(e.target.value)}
						type="password"
						className="border-1 border-accent px-2 py-2.5 rounded focus:outline-none focus:ring-2 focus:ring-accent"
						placeholder="•••••••••"
						value={password}
						id="password"
					/>
				</form>
				<button
					className="bg-accent text-lg md:text-xl text-secondary py-3 px-8 rounded-xl hover:bg-accent/90 hover:scale-105 transition-all duration-300 focus:outline-none focus:ring-2 focus:ring-accent focus:ring-offset-2"
					onClick={handlerSignIn}
				>
					Sign in
				</button>
			</div>
		</div>
	);
};
