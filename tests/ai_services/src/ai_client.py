"""
AI Client module with various AI service integrations
"""
import os
from typing import Optional, List, Dict, Any

# OpenAI imports and usage
from openai import OpenAI
import openai

# Anthropic Claude imports
from anthropic import Anthropic

# Hugging Face imports
from transformers import AutoTokenizer, AutoModel, AutoModelForCausalLM, pipeline
from huggingface_hub import HfApi, login

# LangChain imports
from langchain import LLMChain, PromptTemplate
from langchain.llms import OpenAI as LangChainOpenAI
from langchain.chat_models import ChatOpenAI

# Other AI services
import replicate
import cohere
import pinecone
import elevenlabs
import together
import groq
from mistralai.client import MistralClient




class AIServiceManager:
    """Manages multiple AI service integrations"""

    def __init__(self):
        # Initialize OpenAI client
        self.openai_client = OpenAI(
            api_key=os.getenv("OPENAI_API_KEY")
        )

        # Initialize Anthropic client
        self.anthropic_client = Anthropic(
            api_key=os.getenv("ANTHROPIC_API_KEY")
        )

        # Initialize Cohere client
        self.cohere_client = cohere.Client(
            api_key=os.getenv("COHERE_API_KEY")
        )

        # Initialize Pinecone
        pinecone.init(
            api_key=os.getenv("PINECONE_API_KEY"),
            environment=os.getenv("PINECONE_ENVIRONMENT")
        )

        # Initialize Groq client
        self.groq_client = groq.Groq(
            api_key=os.getenv("GROQ_API_KEY")
        )

                # Initialize Mistral client
        self.mistral_client = MistralClient(
            api_key=os.getenv("MISTRAL_API_KEY")
        )

    def generate_with_gpt4(self, prompt: str) -> str:
        """Generate text using GPT-4"""
        response = self.openai_client.chat.completions.create(
            model="gpt-4",
            messages=[{"role": "user", "content": prompt}]
        )
        return response.choices[0].message.content

    def generate_with_claude(self, prompt: str) -> str:
        """Generate text using Claude-3"""
        response = self.anthropic_client.messages.create(
            model="claude-3-sonnet-20240229",
            max_tokens=1000,
            messages=[{"role": "user", "content": prompt}]
        )
        return response.content[0].text

    def generate_with_replicate(self, model: str, input_data: Dict[str, Any]) -> Any:
        """Generate using Replicate models"""
        return replicate.run(model, input=input_data)

    def create_embeddings_with_openai(self, text: str) -> List[float]:
        """Create embeddings using OpenAI"""
        response = self.openai_client.embeddings.create(
            model="text-embedding-ada-002",
            input=text
        )
        return response.data[0].embedding

    def search_with_pinecone(self, query_vector: List[float], top_k: int = 10) -> Dict:
        """Search vectors in Pinecone"""
        index = pinecone.Index("my-index")
        return index.query(
            vector=query_vector,
            top_k=top_k,
            include_metadata=True
        )

    def generate_speech_with_elevenlabs(self, text: str, voice_id: str) -> bytes:
        """Generate speech using ElevenLabs"""
        return elevenlabs.generate(
            text=text,
            voice=voice_id,
            api_key=os.getenv("ELEVENLABS_API_KEY")
        )

    def load_huggingface_model(self, model_name: str):
        """Load a Hugging Face model"""
        tokenizer = AutoTokenizer.from_pretrained(model_name)
        model = AutoModelForCausalLM.from_pretrained(model_name)
        return tokenizer, model

    def create_langchain_pipeline(self) -> LLMChain:
        """Create a LangChain pipeline"""
        llm = ChatOpenAI(
            temperature=0.7,
            openai_api_key=os.getenv("OPENAI_API_KEY")
        )

        prompt = PromptTemplate(
            input_variables=["question"],
            template="Answer the following question: {question}"
        )

                return LLMChain(llm=llm, prompt=prompt)