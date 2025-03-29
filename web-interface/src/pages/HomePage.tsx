import React, {useState} from "react";
import { Tabs, Tab, Card, CardBody, Tooltip } from "@heroui/react";
import { useNavigate } from "react-router-dom";
import UploadComponent from "../components/UploadComponent";
import ShareComponent from "../components/ShareComponent";
import { FaLink } from "react-icons/fa6";
import { LuUpload } from "react-icons/lu";
import { IoMdInformationCircle } from "react-icons/io"

type Key = string | number;

const HomePage: React.FC = () => {
    const [selectedTab, setSelectedTab] = useState<Key>("upload")
    return (
        <div className="flex w-full h-full justify-center items-center">
            {/* Rectangle */}
            <Card className="max-w-full w-[440px]  bg-wprimary -mt-10">
                <CardBody className="justify-center items-center">
                    <Tabs
                        fullWidth
                        selectedKey={selectedTab}
                        size="sm"
                        aria-label="Tabs colors" 
                        radius="sm" 
                        variant="solid"
                        onSelectionChange={(key) => setSelectedTab(key)}
                        className=""
                        classNames={{
                            base: "bg-wprimary", 
                            tabList: "border border-gray-200 bg-wsecondary p-1 w-[70%] h-9 mx-auto mt-1",
                            tab: "px-4 py-2 text-gray-500 font-normal cursor-pointer transition-all hover:bg-gray-100 data-[selected]:shadow-[0_4px_16px_0_rgba(0,0,0,0.12)] data-[selected]:font-medium data-[selected]:text-gray-800 text-[13px] h-[28px]"
                          }}
                    >
                        <Tab key="upload" title={
                            <div className="flex flex-row justify-center items-center ">
                                <LuUpload size={15} className='mr-2'/>
                                <span className="">Upload</span>
                            </div>
                        } className="w-full">
                            <div className="-my-1 -mb-2">
                                    <div className="flex flex-col justify-center items-center mt-1 mb-1">
                                        <div className="flex flex-row justify-center items-center">
                                            <span className="font-brsonoma font-normal text-gray-950 text-md mr-1">
                                                Upload to <span className="bg-secondary-gradient bg-clip-text text-transparent font-bold">Admin's</span> Folder
                                            </span>
                                            <Tooltip  placement="bottom" content="Uploads here will appear in your computer's folder you configured during installation">
                                                <IoMdInformationCircle className='mb-1'/>
                                            </Tooltip>
                                        </div>
                                        <span className="font-brsonoma font-light text-gray-500 text-xs">
                                            Send files to your computer or invite friends
                                        </span>
                                    </div>
                                <UploadComponent/>
                            </div>
                        </Tab>
                        <Tab key="share" title={
                            <div className="flex items-center ">
                                <FaLink size={16} className='mr-2'/>
                                <span className="">Share</span>
                            </div>
                            } className="w-full">
                            <ShareComponent/>
                        </Tab>
                    </Tabs>
                </CardBody>
            </Card>
        </div>
    );
};

export default HomePage;
