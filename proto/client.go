package pbclient

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/erikbayerlein/mult-protocol-clients/internal/auth"
	pb "github.com/erikbayerlein/mult-protocol-clients/internal/pb"
	"github.com/erikbayerlein/mult-protocol-clients/internal/tcp"
	"google.golang.org/protobuf/proto"
)

type ProtobufClient struct {
	Host string
	Port int
}

func (pc *ProtobufClient) Login(studentId int) error {
	req := &pb.Requisicao{
		Conteudo: &pb.Requisicao_Auth{
			Auth: &pb.Auth{
				AlunoId:   fmt.Sprintf("%d", studentId),
				Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
			},
		},
	}

	payload, err := proto.Marshal(req)
	if err != nil {
		return fmt.Errorf("serialization error: %w", err)
	}

	respBytes, err := tcp.RequestBytes(payload, pc.Host, pc.Port)
	if err != nil {
		return fmt.Errorf("tcp auth error: %w", err)
	}

	var resp pb.Resposta
	if err := proto.Unmarshal(respBytes, &resp); err != nil {
		return fmt.Errorf("decode auth error: %w", err)
	}

	op := resp.GetOperacao()
	if op == nil {
		return fmt.Errorf("invalid auth response")
	}

	token, ok := op.Resultado["token"]
	if !ok || token == "" {
		return fmt.Errorf("token not identified")
	}

	fmt.Println("Received token:", token)

	if err := auth.SaveToken(auth.TokenRecord{StudentId: studentId, Token: token}); err != nil {
		return fmt.Errorf("error saving token: %w", err)
	}
	return nil
}

func (pc *ProtobufClient) Logout(token string) error {
	resp, err := pc.doOperation("logout", token, nil)
	if err != nil {
		return err
	}
	fmt.Println("Logout:", resp)
	return nil
}

func (pc *ProtobufClient) Run(op string, args []string) error {
	rec, err := auth.RequireLogin()
	if err != nil {
		return err
	}
	token := rec.Token

	switch op {
	case "echo":
		if len(args) < 1 {
			return fmt.Errorf("echo needs a message")
		}
		resp, err := pc.doOperation("echo", token, map[string]string{
			"mensagem": strings.Join(args, " "),
		})
		fmt.Println("→", resp)
		return err

	case "sum":
		if len(args) < 1 {
			return fmt.Errorf("sum needs a list of nums, ex: 1,2,3")
		}
		resp, err := pc.doOperation("soma", token, map[string]string{
			"numeros": args[0],
		})
		fmt.Println("→", resp)
		return err

	case "timestamp":
		resp, err := pc.doOperation("timestamp", token, nil)
		fmt.Println("→", resp)
		return err

	case "status":
		resp, err := pc.doOperation("status", token, map[string]string{
			"detalhado": "true",
		})
		fmt.Println("→", resp)
		return err

	case "history", "historico":
		limit := "10"
		if len(args) >= 1 {
			limit = args[0]
		}
		resp, err := pc.doOperation("historico", token, map[string]string{
			"limite": limit,
		})
		fmt.Println("→", resp)
		return err

	default:
		return fmt.Errorf("unknown operation: %s", op)
	}
}

func (pc *ProtobufClient) doOperation(nomeOperacao, token string, params map[string]string) (string, error) {
	req := &pb.Requisicao{
		Conteudo: &pb.Requisicao_Operacao{
			Operacao: &pb.Operacao{
				Token:        token,
				NomeOperacao: nomeOperacao,
				Parametros:   params,
				Timestamp:    time.Now().UTC().Format(time.RFC3339Nano),
			},
		},
	}

	payload, err := proto.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("serialization error: %w", err)
	}

	respBytes, err := tcp.RequestBytes(payload, pc.Host, pc.Port)
	if err != nil {
		return "", fmt.Errorf("tcp error: %w", err)
	}

	var resp pb.Resposta
	if err := proto.Unmarshal(respBytes, &resp); err != nil {
		return "", fmt.Errorf("decode response error: %w", err)
	}

	op := resp.GetOperacao()
	if op == nil && nomeOperacao != "logout" {
		return "", fmt.Errorf("invalid response")
	}

	return formatResultado(nomeOperacao, op), nil
}

func formatResultado(cmd string, op *pb.OperacaoResponse) string {
	if op == nil {
		return "ok"
	}

	var b strings.Builder
	b.WriteString("ok {\n")
	b.WriteString(fmt.Sprintf(`  comando: "%s"`+"\n", strings.ToUpper(cmd)))

	keys := make([]string, 0, len(op.Resultado))
	for k := range op.Resultado {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := op.Resultado[k]
		b.WriteString("  dados {\n")
		b.WriteString(fmt.Sprintf("    key: %q\n", k))
		b.WriteString(fmt.Sprintf("    value: %q\n", v))
		b.WriteString("  }\n")
	}

	if op.Timestamp != "" {
		b.WriteString(fmt.Sprintf("  timestamp: %q\n", op.Timestamp))
	}
	b.WriteString("}")

	return b.String()
}
