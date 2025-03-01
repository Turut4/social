package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"

	"github.com/turut4/social/internal/store"
)

var usernames = []string{
	"Bob", "Carlos", "Pietra", "Julia", "Miguel", "Sofia", "Lucas", "Ana",
	"Maria", "Jhon", "Laura", "Gabriel", "Isabella", "Rafael", "Amanda",
	"Diego", "Felipe", "Helena", "Eduardo", "Mariana", "Bruno", "Camila",
	"Ricardo", "Larissa", "Thiago", "Alice", "Daniel", "Luiza", "Rodrigo",
	"Paul", "Clara", "Andre", "Caroline", "Gustavo", "Isabela", "Victor",
	"Alex", "Marcel", "Bianca", "Matheus", "Renata",
}

var titles = []string{
	"Digital Revolution", "Healthy Living", "Japan Journey", "Vegan Recipes", "Study Tips", "Tech News 2024", "Gaming World", "Coffee Culture",
	"Urban Photography", "Indie Music Guide", "Modern Cinema", "Extreme Sports", "Environment Now", "Modern Art Today", "Pets at Home", "Books of the Month",
	"Daily Mindfulness", "Career Development", "Fashion Trends", "Home Decor", "Science Today", "Ancient History", "Personal Finance", "Space Exploration",
	"Basic Gardening", "City Cycling", "Yoga for All", "Language Learning", "Tech Universe",
}

var contents = []string{
	"Exploring the latest technology trends and their impact on our daily lives", "Tips and tricks for maintaining a healthy lifestyle in the modern world",
	"A deep dive into Japanese culture, traditions and must-visit locations", "Delicious plant-based recipes that will revolutionize your cooking",
	"Expert strategies to improve your learning efficiency and academic success", "Breaking down the most significant tech developments of 2024",
	"Reviews and insights about the newest games and gaming industry trends", "Discovering coffee cultures around the world and brewing techniques",
	"Capturing the essence of city life through the lens of street photography", "Spotlight on emerging artists and the evolving indie music scene",
	"Analysis of contemporary cinema and its influence on modern culture", "Adrenaline-pumping adventures and extreme sports experiences",
	"Critical environmental issues and solutions for a sustainable future", "Contemporary art movements and their impact on society",
	"Essential tips for pet care and creating a pet-friendly home environment", "Must-read books and literary analysis for the avid reader",
	"Practical mindfulness exercises for daily stress management", "Strategic advice for career growth in the modern workplace",
	"Latest trends from global fashion weeks and street style inspiration", "Creative ideas for transforming your living space",
	"Breakthrough scientific discoveries and their real-world applications", "Fascinating stories from ancient civilizations and their legacy",
	"Smart investment strategies and personal finance management tips", "Latest discoveries in astronomy and space exploration missions",
	"Beginner-friendly gardening tips for urban and suburban spaces", "Urban cycling guides and bicycle maintenance essentials",
	"Accessible yoga practices for practitioners of all levels", "Effective methods for mastering new languages",
	"Deep dives into emerging technologies and digital transformation", "Inspiring travel narratives and destination guides",
	"Advanced cooking techniques and international cuisine exploration", "Modern architectural trends and innovative design concepts",
	"Digital marketing strategies and social media optimization tips", "Understanding mental health and wellness practices",
	"Comprehensive fitness programs and nutrition guidance", "Latest updates from the world of comics, games, and tech",
	"Success stories and practical advice for aspiring entrepreneurs", "Innovative approaches to environmental conservation",
	"Updates and theories about the expanding Marvel cinematic universe", "Expert travel tips and hidden gem destinations",
	"Methods for maximizing productivity in work and personal life", "Step-by-step guides for creative home projects",
	"Understanding wine varieties and food pairing fundamentals", "Exploring human behavior and psychological concepts",
	"Latest automotive innovations and industry developments", "Analyzing trends and phenomena in popular culture",
}

var tags = []string{
	"Technology", "Health", "Travel", "Food", "Education", "Gaming", "Photography", "Culture",
	"Music", "Entertainment", "Sports", "Environment", "Art", "Pets", "Literature", "Lifestyle",
	"Career", "Fashion", "Science", "History", "Finance", "Space", "Gardening", "Fitness",
	"Yoga", "Languages", "Business", "Cooking", "Architecture", "Marketing", "Wellness", "Geek",
	"Innovation", "Nature", "Movies", "Adventure", "DIY", "Psychology", "Automotive", "Social Media",
	"Development", "Programming", "Design", "Creativity", "Leadership", "Motivation",
}

var comments = []string{
	"Great article! Really helped me understand the topic better.", "Thanks for sharing this valuable information!", "This is exactly what I was looking for.", "Very well written and informative.",
	"I've learned something new today, thank you!", "Looking forward to more content like this!", "Interesting perspective on this subject.", "Could you make a follow-up post about this?",
	"Brilliant analysis of the topic.", "The examples really clarified the concept.", "Sharing this with my colleagues!", "This changed my view on the subject.",
	"Excellent research and presentation.", "Your writing style is so engaging!", "Never thought about it this way before.", "This deserves more attention!",
	"Clear and concise explanation.", "Would love to see more content like this.", "The tips you shared are really practical.", "This helped me solve my problem.",
	"Really comprehensive coverage of the topic.", "Such an interesting take on this!", "You've covered all the important aspects.", "This is going in my bookmarks.",
	"Fantastic breakdown of complex ideas.", "The insights here are invaluable.", "Great work on explaining this clearly.", "This is exactly what the industry needs.",
	"Love how you structured this information.", "Can't wait to implement these ideas.", "You've made this topic so accessible.", "Perfect timing, I needed this info!",
	"Amazing resource for beginners.", "Your expertise really shows here.", "This answered all my questions.", "Well-researched and informative.",
	"Finally, someone explained this properly!", "This is gold! Saving for later.", "Outstanding content as always.", "Such a helpful resource!",
}

func Seed(store store.Storage, db *sql.DB) {
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return
	}

	users := generateUsers(300)
	for _, user := range users {
		if err := store.Users.Create(ctx, tx, user); err != nil {
			_ = tx.Rollback()
			log.Println("Error creating user: ", err)
			return
		}
	}

	tx.Commit()

	posts := generatePosts(500, users)
	for _, post := range posts {
		if err := store.Posts.Create(ctx, post); err != nil {
			log.Println("Error creating post: ", err)
			return
		}
	}

	comments := generateComments(1000, users, posts)
	for _, comment := range comments {
		if err := store.Comments.Create(ctx, comment); err != nil {
			log.Println("Error creating comment: ", err)
			return
		}
	}

	log.Println("seed has been completed")
}

func generateUsers(num int) []*store.User {
	users := make([]*store.User, num)

	for i := 0; i < num; i++ {
		users[i] = &store.User{
			Username: usernames[i%len(usernames)] + fmt.Sprintf("%d", i),
			Email:    usernames[i%len(usernames)] + fmt.Sprintf("%d", i) + "@example.com",
			RoleID: 1,
		}
	}

	return users
}

func generatePosts(num int, users []*store.User) []*store.Post {
	posts := make([]*store.Post, num)
	for i := 0; i < num; i++ {
		numTags := rand.Intn(5)
		postTags := []string{}
		for j := 0; j < numTags; j++ {
			postTags = append(postTags, tags[rand.Intn(len(tags))])
		}

		posts[i] = &store.Post{
			UserID:  users[rand.Intn(len(users))].ID,
			Title:   titles[rand.Intn(len(titles))],
			Content: contents[rand.Intn(len(contents))],
			Tags:    postTags,
		}
	}

	return posts
}

func generateComments(num int, users []*store.User, posts []*store.Post) []*store.Comment {
	cms := make([]*store.Comment, num)

	for i := 0; i < num; i++ {
		cms[i] = &store.Comment{
			UserID:  users[rand.Intn(len(users))].ID,
			PostID:  posts[rand.Intn(len(posts))].ID,
			Content: comments[rand.Intn(len(comments))],
		}
	}

	return cms
}
