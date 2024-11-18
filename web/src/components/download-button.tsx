import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faArrowCircleDown } from '@fortawesome/free-solid-svg-icons'

import { Button } from "@/components/ui/button"
import {
    Popover,
    PopoverContent,
    PopoverTrigger,
} from "@/components/ui/popover"

export default function DownloadButton() {
    return (
        <Popover>

            <PopoverTrigger >
                <Button size={"lg"} className="rounded-full pointer-events-none" asChild>
                    <a href={"/download"}>
                        <FontAwesomeIcon icon={faArrowCircleDown} className='text-lg' />
                        <span className="ml-2 text-lg">Скачать</span>
                    </a>
                </Button>
            </PopoverTrigger>

            <PopoverContent>
                <p>Скоро релиз</p>
            </PopoverContent>

        </Popover>
    )
}