import { useNavigate, useParams } from "react-router-dom";
import { API_URL } from "./App";

export const ConfimationPage = () => {
	const { token = "" } = useParams();
	const redirect = useNavigate();

	const handleConfirm = async () => {
		const response = await fetch(`${API_URL}/users/activate/${token}`, {
			method: "PUT",
		});

		console.log(response.headers);
		if (response.ok) {
			redirect("/");
		} else {
			alert("failed to confirm");
		}
	};
	return (
		<div className="flex bg-background h-screen items-center justify-center p-5">
			<div className="bg-secondary h-[400px] w-[320px] md:h-[450px] md:w-[400px] flex flex-col items-center justify-center rounded-3xl gap-8 text-primary shadow-lg border border-border">
				<h1 className="text-4xl md:text-5xl text-accent font-bold mb-4">
					Confirmation
				</h1>
				<p className="text-lg md:text-xl p-6 text-center text-gray-600 leading-relaxed">
					Confirm your email to start a new experience with{" "}
					<span className="text-accent font-semibold">Chime</span>
				</p>
				<button
					className="bg-accent text-lg md:text-xl text-secondary py-3 px-8 rounded-xl hover:bg-accent/90 hover:scale-105 transition-all duration-300 focus:outline-none focus:ring-2 focus:ring-accent focus:ring-offset-2"
					onClick={handleConfirm}
				>
					Click to confirm
				</button>
			</div>
		</div>
	);
};
