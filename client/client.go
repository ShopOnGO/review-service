package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/ShopOnGO/review-proto/pkg/service"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	reviewClient := pb.NewReviewServiceClient(conn)
	questionClient := pb.NewQuestionServiceClient(conn)

	DeleteAllReviews(reviewClient)
	testReviewService(reviewClient)

	DeleteAllQuestions(questionClient)
	testQuestionService(questionClient)
}

func testReviewService(client pb.ReviewServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	createResp, err := client.AddReview(ctx, &pb.AddReviewRequest{
		ProductVariantId: 1,
		UserId:           1,
		Rating:           5,
		Comment:          "Excellent product!",
	})
	if err != nil {
		log.Fatalf("error creating review: %v", err)
	}
	fmt.Printf("✅ Review created: ID=%d, Rating=%d, Comment=%s\n",
		createResp.Review.Model.Id, createResp.Review.Rating, createResp.Review.Comment)

	getResp, err := client.GetReviews(ctx, &pb.GetReviewsRequest{
		ProductVariantId: 1,
	})
	if err != nil {
		log.Fatalf("error getting reviews: %v", err)
	}
	fmt.Println("✅ Reviews for product variant 1:")
	for _, r := range getResp.Reviews {
		fmt.Printf("   - ID=%d, Rating=%d, Comment=%s\n", r.Model.Id, r.Rating, r.Comment)
	}

	updateResp, err := client.UpdateReview(ctx, &pb.UpdateReviewRequest{
		ReviewId: createResp.Review.Model.Id,
		Rating:   4,
		Comment:  "Good product, but could be better.",
	})
	if err != nil {
		log.Fatalf("error updating review: %v", err)
	}
	fmt.Printf("✅ Review updated: Success=%v, Message=%s\n", updateResp.Success, updateResp.Message)

	avgResp, err := client.GetAverageRating(ctx, &pb.GetAverageRatingRequest{
		ProductVariantId: 1,
	})
	if err != nil {
		log.Fatalf("error getting average rating: %v", err)
	}
	fmt.Printf("✅ Average rating for product variant 1: %f\n", avgResp.AverageRating)

	delResp, err := client.DeleteReview(ctx, &pb.DeleteReviewRequest{
		ReviewId: createResp.Review.Model.Id,
	})
	if err != nil {
		log.Fatalf("error deleting review: %v", err)
	}
	fmt.Printf("✅ Review deleted: Success=%v, Message=%s\n", delResp.Success, delResp.Message)
}

func DeleteAllReviews(client pb.ReviewServiceClient) {
	fmt.Println("⚠ DeleteAllReviews: function not implemented, skipping.")
}

func testQuestionService(client pb.QuestionServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	addResp, err := client.AddQuestion(ctx, &pb.AddQuestionRequest{
		ProductVariantId: 1,
		UserId:           1,
		QuestionText:     "Is this product waterproof?",
	})
	if err != nil {
		log.Fatalf("error adding question: %v", err)
	}
	fmt.Printf("✅ Question added: ID=%d, Text=%s\n", addResp.Question.Model.Id, addResp.Question.QuestionText)

	getResp, err := client.GetQuestions(ctx, &pb.GetQuestionsRequest{
		ProductVariantId: 1,
	})
	if err != nil {
		log.Fatalf("error getting questions: %v", err)
	}
	fmt.Println("✅ Questions for product variant 1:")
	for _, q := range getResp.Questions {
		fmt.Printf("   - ID=%d, Text=%s, Answer=%s\n", q.Model.Id, q.QuestionText, q.AnswerText)
	}

	answerResp, err := client.AnswerQuestion(ctx, &pb.AnswerQuestionRequest{
		QuestionId: addResp.Question.Model.Id,
		AnswerText: "Yes, it is fully waterproof.",
	})
	if err != nil {
		log.Fatalf("error answering question: %v", err)
	}
	fmt.Printf("✅ Question answered: Success=%v, Message=%s\n", answerResp.Success, answerResp.Message)

	delResp, err := client.DeleteQuestion(ctx, &pb.DeleteQuestionRequest{
		QuestionId: addResp.Question.Model.Id,
	})
	if err != nil {
		log.Fatalf("error deleting question: %v", err)
	}
	fmt.Printf("✅ Question deleted: Success=%v, Message=%s\n", delResp.Success, delResp.Message)
}

func DeleteAllQuestions(client pb.QuestionServiceClient) {
	fmt.Println("⚠ DeleteAllQuestions: function not implemented, skipping.")
}
