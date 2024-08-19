from typing import Tuple
from pathlib import Path

from yaml import Mark


from .niclib import rand_string


import subprocess

import requests
import logging

from typing import List, Optional, Any

from warnings import warn

import os

import aiohttp
import asyncio

from pydantic import BaseModel

from constants import (
    OS_TMPDIR,
    DATALAB_API_KEY,
    MARKER_ENDPOINT_URL,
)


global_marker_server_urls = ["https://marker.kessler.xyz"]


class GPUComputeEndpoint:
    def __init__(
        self,
        logger: Any,
        marker_endpoint_url: str = MARKER_ENDPOINT_URL,
        legacy_endpoint_url: str = "https://depricated-url.com",
        datalab_api_key: str = DATALAB_API_KEY,
    ):
        self.logger = logger
        self.marker_endpoint_url = marker_endpoint_url
        self.endpoint_url = legacy_endpoint_url
        self.datalab_api_key = datalab_api_key

    async def pull_marker_endpoint_for_response(
        self, request_check_url: str, max_polls: int, poll_wait: int
    ) -> str:
        async with aiohttp.ClientSession() as session:
            for polls in range(max_polls):
                try:
                    await asyncio.sleep(poll_wait)
                    async with session.get(request_check_url) as poll_response:
                        poll_data = await poll_response.json()
                        # self.logger.info(poll_data)
                        if poll_data["status"] == "complete":
                            self.logger.info(f"Processed document after {polls} polls.")
                            return poll_data["markdown"]
                        if poll_data["status"] == "error":
                            e = poll_data["error"]
                            self.logger.error(
                                f"Pdf server encountered an error after {polls} : {e}"
                            )
                            raise Exception(
                                f"Pdf server encountered an error after {polls} : {e}"
                            )
                        if poll_data["status"] != "processing":
                            raise ValueError(
                                f"PDF Processing Failed. Status was unrecognized {poll_data['status']} after {polls} polls."
                            )
                except Exception as e:
                    raise e
            raise TimeoutError("Polling for marker API result timed out")

    async def transcribe_pdf_s3_uri(
        self, s3_uri: str, external_process: bool = False, priority: bool = True
    ) -> str:
        if external_process:
            # TODO : Make it so that it downloads the s3_uri onto local then uploads it to external process.
            raise Exception(
                "s3 uploads not supported with external process equaling true."
            )
        else:
            base_url = global_marker_server_urls[0]
            if priority:
                query_str = "?priority=true"
            else:
                query_str = "?priority=false"
            marker_url_endpoint = (
                base_url + "/api/v1/marker/direct_s3_url_upload" + query_str
            )

            data = {"s3_url": s3_uri}
            # data = {"langs": "en", "force_ocr": "false", "paginate": "true"}
            async with aiohttp.ClientSession() as session:
                async with session.post(marker_url_endpoint, json=data) as response:
                    response_data = await response.json()
                    # await the json if async
                    request_check_url_leaf = response_data.get("request_check_url_leaf")

                    if request_check_url_leaf is None:
                        raise Exception(
                            "Failed to get request_check_url from marker API response"
                        )
                    request_check_url = base_url + request_check_url_leaf
                    self.logger.info(
                        f"Got response from marker server, polling to see when file is finished processing at url: {request_check_url}"
                    )
            return await self.pull_marker_endpoint_for_response(
                request_check_url=request_check_url,
                max_polls=200,
                poll_wait=3 + 57 * int(not priority),
            )

    async def transcribe_pdf_filepath(
        self,
        filepath: Path,
        external_process: bool = False,
        priority=True,
    ) -> str:
        if external_process:
            url = "https://www.datalab.to/api/v1/marker"
            self.logger.info(
                "Calling datalab api with key beginning with"
                + self.datalab_api_key[0 : (len(self.datalab_api_key) // 5)]
            )
            headers = {"X-Api-Key": self.datalab_api_key}

            with open(filepath, "rb") as file:
                files = {
                    "file": (filepath.name + ".pdf", file, "application/pdf"),
                    "paginate": (None, True),
                }
                # data = {"langs": "en", "force_ocr": "false", "paginate": "true"}
                with requests.post(url, files=files, headers=headers) as response:
                    response.raise_for_status()
                    # await the json if async
                    data = response.json()
                    request_check_url = data.get("request_check_url")

                    if request_check_url is None:
                        raise Exception(
                            "Failed to get request_check_url from marker API response"
                        )
                    self.logger.info(
                        "Got response from marker server, polling to see when file is finished processing."
                    )
                    return await self.pull_marker_endpoint_for_response(
                        request_check_url=request_check_url,
                        max_polls=200,
                        poll_wait=3 + 57 * int(not priority),
                        server=MarkerServer("void", 0),
                    )
        else:
            server = get_least_connection()
            base_url = server.url
            if priority:
                query_str = "?priority=true"
            else:
                query_str = "?priority=false"
            marker_url_endpoint = base_url + "/api/v1/marker" + query_str

            with open(filepath, "rb") as file:
                files = {
                    "file": (filepath.name + ".pdf", file, "application/pdf"),
                    # "paginate": (None, True),
                }
                # data = {"langs": "en", "force_ocr": "false", "paginate": "true"}
                with requests.post(marker_url_endpoint, files=files) as response:
                    response.raise_for_status()
                    # await the json if async
                    data = response.json()
                    request_check_url_leaf = data.get("request_check_url_leaf")

                    if request_check_url_leaf is None:
                        raise Exception(
                            "Failed to get request_check_url from marker API response"
                        )
                    request_check_url = base_url + request_check_url_leaf
                    self.logger.info(
                        f"Got response from marker server, polling to see when file is finished processing at url: {request_check_url}"
                    )
                    return await self.pull_marker_endpoint_for_response(
                        request_check_url=request_check_url,
                        max_polls=200,
                        poll_wait=3 + 57 * int(not priority),
                        server=server,
                    )

    def audio_to_text_raw(
        self, filepath: Path, source_lang: str, target_lang: str, file_type: str
    ) -> str:
        # The API endpoint you will be hitting
        url = f"{self.endpoint_url}/v0/multimodal_asr/whisper-latest"
        # Open the file in binary mode
        with filepath.open("rb") as file:
            # Define the multipart/form-data payload
            files = {"file": (filepath.name, file, "application/octet-stream")}
            jsonpayload = {
                "source_lang": source_lang,
                "target_lang": target_lang,
            }
            # Make the POST request with files
            response = requests.post(url, files=files, json=jsonpayload)
            # Raise an exception if the request was unsuccessful
            response.raise_for_status()

        # Parse the JSON response
        response_json = response.json()

        # Extract the translated text from the JSON response
        translated_text = response_json["response"]
        return translated_text

    def audio_to_text(
        self,
        filepath: Path,
        source_lang: str,
        target_lang: str,
        file_type: str,
        bitrate: int = 15000,
        tmpdir: Path = OS_TMPDIR,
    ) -> str:
        downsampled = downsample_audio(filepath, file_type, bitrate, tmpdir)
        if downsampled == filepath:
            return self.audio_to_text_raw(filepath, source_lang, target_lang, file_type)
        return self.audio_to_text_raw(downsampled, source_lang, target_lang, "opus")

    def translate_text(
        self, doctext: str, source_lang: Optional[str], target_lang: str
    ) -> str:
        url = f"{self.endpoint_url}/v0/translation/google-translate"
        payload = {
            "text": doctext,
            "source_lang": source_lang,
            "target_lang": target_lang,
        }
        response = requests.post(url, json=payload)
        response.raise_for_status()
        text = response.json().get("text", [])
        return text
