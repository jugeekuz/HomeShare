import React from "react"
import AppRoutes from "./AppRoutes.tsx"
import {HeroUIProvider} from "@heroui/react";

const App : React.FC = () => {
	return (
        <HeroUIProvider>
		<div
			className="relative items-center w-[100dvw] h-[100dvh] "
			style={{
					backgroundImage: "url('src/assets/img/abstract-bg.png')",
					backgroundSize: "cover",
					backgroundPosition: "center",
					backgroundRepeat: "no-repeat"
			}}
		>
			{/* Gradient Overlay */}
			<div
				className="z-[1] absolute top-0 left-0 right-0 bottom-0"
					style={{
							background: "linear-gradient(to bottom, rgba(0, 0, 0, 0.6), rgba(255, 255, 255, 0))", 
					}}
			/>

			<div className="relative z-[2] h-full w-full flex flex-col items-center justify-center">
                    <AppRoutes/>
			</div>
		</div>
        </HeroUIProvider>
	)
}
export default App
