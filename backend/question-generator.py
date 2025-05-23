import os
import json
import requests
from typing import List, Dict, Any, Optional
import openai
from datetime import datetime

# Initialize OpenAI client with API key from environment
client = openai.OpenAI(api_key=os.getenv("OPENAI_API_KEY"))

class QuestionGenerator:
    def __init__(self, api_url: str, headers: Optional[Dict[str, str]] = None):
        self.api_url = api_url
        self.headers = headers or {
            "Content-Type": "application/json"
        }

    def generate_questions(
        self,
        category: str,
        num_questions: int = 5,
        num_fillers: int = 3,
    ) -> List[Dict[str, Any]]:
        """
        Generate multiple-choice questions for a specific category.
        
        Args:
            category: The category of questions (e.g., "Science", "History")
            num_questions: Number of questions to generate (1-10)
            num_fillers: Number of filler answers per question (3-5)

        Returns:
            List of question dictionaries with text, answer, and filler_answers
        """
        # Input validation
        num_questions = max(1, min(10, num_questions))  # Clamp between 1-10
        num_fillers = max(3, min(5, num_fillers))  # Clamp between 3-5
        
        prompt = f"""Generate {num_questions} questions about {category}.
For each question, provide:
1. A clear, concise question text
2. A single correct answer (factually accurate)
3. {num_fillers} incorrect but plausible filler answers

Important guidelines:
- Questions should be factually accurate and verifiable
- Avoid opinion-based questions
- Ensure answers are not too similar to each other
- Keep questions and answers concise
- Format filler_answers as a JSON array of strings

Format the response as a JSON array of objects with this exact structure:
[
  {{
    "text": "What is the capital of France?",
    "answer": "Paris",
    "filler_answers": ["London", "Berlin", "Madrid"]
  }}
]

ONLY return the JSON array, no other text or markdown formatting."""

        try:
            response = client.chat.completions.create(
                model="gpt-4-turbo",
                messages=[
                    {"role": "system", "content": "You are a helpful assistant that creates high-quality, factually accurate quiz questions."},
                    {"role": "user", "content": prompt}
                ],
                temperature=0.7,
                max_tokens=2000
            )
            
            # Extract and clean the JSON response
            content = response.choices[0].message.content
            content = content.strip().strip("```json").strip("```").strip()
            questions = json.loads(content)
            
            # Add category to each question
            for q in questions:
                q["category"] = category
                
            return questions
            
        except json.JSONDecodeError as e:
            print(f"Error parsing JSON response: {e}")
            print("Raw response:", content if 'content' in locals() else 'No content')
            return []
        except Exception as e:
            print(f"Error generating questions: {e}")
            return []

    def upload_questions(self, questions: List[Dict[str, Any]]) -> bool:
        """
        Upload generated questions to the backend API.
        
        Args:
            questions: List of question dictionaries
            
        Returns:
            bool: True if upload was successful, False otherwise
        """
        if not questions:
            print("No questions to upload")
            return False
            
        # Ensure questions match the expected format
        formatted_questions = []
        for q in questions:
            formatted_questions.append({
                "text": q["text"],
                "answer": q["answer"],
                "category": q["category"],
                "filler_answers": q["filler_answers"]
            })
            
        payload = {"questions": formatted_questions}
        
        try:
            print(f"Sending request to {self.api_url}")
            print("Payload:", json.dumps(payload, indent=2))
            
            response = requests.post(
                self.api_url,
                json=payload,
                headers=self.headers,
                timeout=30  # 30 second timeout
            )
            
            response.raise_for_status()
            print(f"Successfully uploaded {len(questions)} questions")
            return True
            
        except requests.exceptions.RequestException as e:
            print(f"Error uploading questions: {e}")
            if hasattr(e, 'response') and e.response is not None:
                print(f"Status code: {e.response.status_code}")
                print("Response:", e.response.text)
            return False

def main():
    # Configuration
    API_URL = "http://localhost:8080/api/games/questions/bulk"  # Updated endpoint
    CATEGORIES = ["Science", "History", "Geography", "Sports", "Entertainment"]
    QUESTIONS_PER_CATEGORY = 5
    
    # Initialize generator
    generator = QuestionGenerator(API_URL)
    
    # Process each category
    for category in CATEGORIES:
        print(f"\n{'='*50}")
        print(f"Processing category: {category}")
        print(f"{'='*50}")
        
        # Generate questions
        questions = generator.generate_questions(
            category=category,
            num_questions=QUESTIONS_PER_CATEGORY,
            num_fillers=3
        )
        
        if not questions:
            print(f"Failed to generate questions for {category}")
            continue
            
        print(f"Generated {len(questions)} questions for {category}")
        
        # Upload questions
        if generator.upload_questions(questions):
            print(f"Successfully processed {len(questions)} questions for {category}")
        else:
            print(f"Failed to upload questions for {category}")

if __name__ == "__main__":
    main()
