import React from "react";
import { Divider, Button } from "@heroui/react";
import { useNavigate } from "react-router-dom";
import { FileUploadIcon } from "../components/FileIcons";
const HomePage: React.FC = () => {
    const navigate = useNavigate();
    return (
        <div className="flex h-full w-full justify-center items-center px-1">
            <div className="flex flex-col justify-center items-center w-[25rem] xl:w-[28rem] bg-white rounded-xl shadow-lg -mt-20">
                <div className="flex flex-col items-center justify-center w-full h-64">
                    <div className="flex flex-col items-center justify-center w-full h-full">
                        <div className="flex flex-row w-full h-full p-2">
                            <div className="flex flex-col h-full w-1/2 items-center justify-center">
                            {/* Left Container */}
                                {/* <div className="p-2">
                                    <span className="font-linik text-gray-800 text-xs mt-2">
                                        Upload Files Directly
                                    </span>
                                </div> */}
                                <FileUploadIcon className=""/>
                                <div className="p-2">
                                    <span className="font-linik text-gray-800 text-xs mt-2">
                                        Upload Files Directly
                                    </span>
                                </div>
                                <Button
                                    color="primary"
                                    className="text-sm w-3/4"//Change font
                                    size="md"
                                    radius="sm"
                                    onPress={() => navigate('/upload')}>
                                    Upload Files
                                </Button>
                            </div>
                            <div className="relative flex justify-center items-center h-full">
                                <div className="absolute flex flex-col justify-between items-center h-full">
                                    <Divider orientation="vertical" className="h-2/5"/>
                                    <span className="font-linik text-gray-400 text-xs ">OR</span>
                                    <Divider orientation="vertical" className="h-2/5"/>
                                </div>
                            </div>
                            <div className="flex flex-col h-full w-1/2 items-center justify-center">
                                {/* Right Container */}

                                <Button
                                    color="primary"
                                    className="text-sm w-3/4"//Change font
                                    size="md"
                                    radius="sm"
                                    onPress={() => navigate('/sharing')}>
                                        Share Files
                                </Button>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default HomePage;
