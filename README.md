# Dahaa Game Frontend

A modern web-based implementation of the Dahaa Game using React, TypeScript, and Vite.

## Features

- Token-based authentication
- Responsive design with Tailwind CSS
- Modern UI components with shadcn/ui
- State management with Zustand
- TypeScript for type safety
- PWA support (coming soon)

## Prerequisites

- Node.js 18+ and npm/yarn
- Backend API running (default: http://localhost:8080)

## Setup

1. Clone the repository
2. Install dependencies:

   ```bash
   npm install
   # or
   yarn install
   ```

3. Create a `.env` file in the root directory with the following content:

   ```
   VITE_API_URL=http://localhost:8080
   ```

   Adjust the URL if your backend is running on a different port or host.

4. Start the development server:

   ```bash
   npm run dev
   # or
   yarn dev
   ```

5. Open http://localhost:3000 in your browser

## Available Scripts

- `npm run dev` - Start development server
- `npm run build` - Build for production
- `npm run preview` - Preview production build
- `npm run lint` - Run ESLint

## Project Structure

```
dahaa-frontend/
├── app/                    # Application source code
│   ├── components/        # Reusable components
│   ├── lib/              # Utilities and store
│   ├── pages/            # Page components
│   └── globals.css       # Global styles
├── public/               # Static assets
└── ...config files
```

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

MIT
